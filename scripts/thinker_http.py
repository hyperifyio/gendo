#!/usr/bin/env python3
"""
thinker_http_txt.py  –  Reasoning loop with *plain-text* I/O
============================================================
Uses raw HTTP (requests) to hit any OpenAI-style /chat/completions endpoint.
"""

import os, sys, time, json, re, requests, logging, textwrap
from concurrent.futures import ThreadPoolExecutor
from typing import List, Dict, Any, Tuple, Optional
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

# ──────────────────────────────────────────────────────────────────────
# 0  Configuration & HTTP helper
# ──────────────────────────────────────────────────────────────────────
BASE_URL = os.getenv("OPENAI_BASE_URL", "http://localhost:18080/v1").rstrip("/")
API_KEY  = os.getenv("OPENAI_API_KEY",  "local-key")
MODEL    = os.getenv("OPENAI_MODEL",    "bitnet")

N_CANDIDATES      = int(os.getenv("N_CANDIDATES",      "5"))
MAX_REFINE_ROUNDS = int(os.getenv("MAX_REFINE_ROUNDS", "1"))
MAX_WORKERS       = int(os.getenv("MAX_WORKERS",       "4"))  # For parallel processing

logging.basicConfig(
    level=os.getenv("LOGLEVEL", "INFO").upper(),
    format="%(asctime)s [%(levelname)s] %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S")
log = logging.getLogger("thinker")

# Create a persistent session with connection pooling
session = requests.Session()
session.headers.update({
    "Content-Type": "application/json",
    "Authorization": f"Bearer {API_KEY}",
})
# Set up connection pooling for both HTTP and HTTPS with retries
adapter = requests.adapters.HTTPAdapter(
    pool_maxsize=MAX_WORKERS,
    pool_connections=MAX_WORKERS,
    max_retries=3
)
session.mount("http://", adapter)
session.mount("https://", adapter)

# ──────────────────────────────────────────────────────────────────────
# 1  Prompt templates (plain-text sections)
# ──────────────────────────────────────────────────────────────────────
MARKUP = textwrap.dedent("""
    THOUGHT:
    <step-by-step mathematical reasoning>

    ANSWER:
    <one-line mathematical conclusion>
""").strip()

GEN_SYS = textwrap.dedent("""
    You are a professional assistant who solves problems step by step.
    You MUST ALWAYS use the EXACT format provided, with THOUGHT and ANSWER sections.
    If you discover the claim is false, your ANSWER line must clearly state that and give the divisor(s) or counter-example.
    Do not add any other sections or text.
""").strip()

GEN_USER = textwrap.dedent("""
    QUESTION:
    {question}

    You MUST use EXACTLY this format, with no other text:

    THOUGHT:
    1. First, analyze what the question is asking
    2. Then, break down the problem into logical steps
    3. Work through each step carefully
    4. Make a clear conclusion based on your analysis

    ANSWER:
    <one-line conclusion>
""").strip()

CRIT_SYS = textwrap.dedent("""
    You are a strict logician who spots errors. Your task is to:
    1. Check for logical contradictions in the reasoning
    2. Verify mathematical or factual accuracy
    3. Ensure each step follows from the previous one
    4. Look for unsupported assumptions
    5. Check if the conclusion matches the reasoning
    
    If you find ANY issues, list them specifically. If there are NO issues, output exactly "NO ISSUE".
""").strip()

CRIT_USER = textwrap.dedent("""\
    TASK: Review the reasoning and list any logical flaws OR output exactly **NO ISSUE**.

    REASONING:
    {thought}

    ANSWER:
    {answer}

    Check specifically for:
    1. Contradictions between steps
    2. Mathematical errors
    3. Logical fallacies
    4. Unsupported claims
    5. Mismatch between reasoning and conclusion
""")

REF_SYS = ("You are the original author. Rewrite your answer fixing every issue "
           "the critic found. Use the SAME THOUGHT/ANSWER template.")
REF_USER = textwrap.dedent("""\
    YOUR PREVIOUS ANSWER:
    {prev_text}

    CRITIC SAYS:
    {critic}
""")

JUDGE_SYS = textwrap.dedent("""
    You are an impartial judge who extracts reusable heuristics from successful reasoning.
    Your task is to:
    1. Choose the single best answer based on:
       - Factual correctness (even if it contradicts the prompt's assumption)
       - Proper THOUGHT/ANSWER format
       - Clear step-by-step reasoning
    2. Extract general rules that made the answer successful
    3. Format rules as reusable heuristics that could apply to similar problems
    4. Focus on reasoning patterns, not specific facts
    5. Make rules abstract enough to be widely applicable
""").strip()

JUDGE_USER = textwrap.dedent("""\
    Below are the candidate answers tagged A, B, C… Choose the single best one
    based on factual correctness and proper reasoning (even if it contradicts the prompt).
    Then extract general heuristics that made the answer successful; begin each rule with "RULE:".

    {candidates}

    FORMAT:
    Best: <letter>
    RULES:
    - RULE: <abstract, reusable heuristic>
    - RULE: <another general principle>
""")

# Pre-build system prompts
SYSTEM_PROMPTS = {
    "gen": {"role": "system", "content": GEN_SYS},
    "crit": {"role": "system", "content": CRIT_SYS},
    "ref": {"role": "system", "content": REF_SYS},
    "judge": {"role": "system", "content": JUDGE_SYS}
}

def chat(msgs: List[Dict[str, str]], temperature: float = 0.7, max_tokens: int = 128) -> str:
    """Make a chat request with retries."""
    prompt_type = "unknown"
    if msgs[0]["content"] == GEN_SYS:
        prompt_type = "generation"
    elif msgs[0]["content"] == CRIT_SYS:
        prompt_type = "critique"
    elif msgs[0]["content"] == REF_SYS:
        prompt_type = "refinement"
    elif msgs[0]["content"] == JUDGE_SYS:
        prompt_type = "judgment"

    log.debug(f"Making chat request with temperature={temperature}, max_tokens={max_tokens}")
    
    # Create local copy of messages
    msgs_local = list(msgs)
    
    # Log full request body
    request_body = {
        "model": MODEL,
        "messages": msgs_local,
        "temperature": temperature,
        "max_tokens": max_tokens,
        "stop": ["\n\nRULE:"]  # Only stop at RULE to allow full THOUGHT/ANSWER
    }
    log.debug("Request body:\n%s", json.dumps(request_body, indent=2))
    
    start_time = time.time()
    try:
        url = f"{BASE_URL}/chat/completions"
        resp = session.post(
            url,
            json=request_body,
            timeout=60  # Increased timeout for longer generations
        )
        resp.raise_for_status()
        response_json = resp.json()
        
        # Log full response body
        log.debug("Response body:\n%s", json.dumps(response_json, indent=2))
        
        result = response_json["choices"][0]["message"]["content"]
        elapsed = time.time() - start_time
        log.info(f"Prompt type '{prompt_type}' took {elapsed:.2f}s")
        return result
    except requests.exceptions.RequestException as e:
        elapsed = time.time() - start_time
        log.error(f"Request failed after {elapsed:.2f}s: {str(e)}")
        raise

# ──────────────────────────────────────────────────────────────────────
# 2  Parsing helpers
# ──────────────────────────────────────────────────────────────────────
# Relaxed regex to match various formats of THOUGHT/ANSWER tags
CAND_RE = re.compile(r"(?:THOUGHT|Thought|thought|REASONING|Reasoning|reasoning):\s*(.*?)\s*(?:ANSWER|Answer|answer|CONCLUSION|Conclusion|conclusion):\s*(.*)", re.S | re.I)

def parse_candidate(text: str) -> Dict[str, str]:
    """Return {'thought': str, 'answer': str, 'raw': str} or None."""
    # Remove backticks, markdown formatting, and normalize whitespace
    text = re.sub(r'[*_`]', '', text.strip())
    
    # Try to find THOUGHT/ANSWER pattern
    m = CAND_RE.search(text)
    if not m:
        # Try to split on double newlines and look for patterns
        parts = text.split('\n\n')
        for i in range(len(parts)-1):
            combined = f"THOUGHT: {parts[i]}\nANSWER: {parts[i+1]}"
            m = CAND_RE.search(combined)
            if m:
                break
    
    if not m:
        log.debug(f"Failed to parse candidate. Text:\n{text.replace('\n', '\\n')}")
        return None
        
    thought = m.group(1).strip()
    answer = m.group(2).strip()
    
    # Normalize whitespace in both sections
    thought = re.sub(r'\s+', ' ', thought)
    answer = re.sub(r'\s+', ' ', answer)
    
    return {"thought": thought,
            "answer": answer,
            "raw": f"THOUGHT:\n{thought}\n\nANSWER:\n{answer}"}

def normalize_critic(text: str) -> str:
    """Normalize critic response for consistent comparison."""
    # Remove all non-alphanumeric characters and convert to uppercase
    return re.sub(r'[^a-zA-Z0-9]', '', text).upper()

# ──────────────────────────────────────────────────────────────────────
# 3  Pipeline stages
# ──────────────────────────────────────────────────────────────────────
def check_semantic_consistency(question: str, answer: str) -> Tuple[bool, str]:
    """Check if the answer's conclusion matches the question's intent."""
    # Check if this is a truth-assessment task
    truth_assessment = any(phrase in question.lower() for phrase in [
        'prove', 'show', 'is it true', 'true or false', 'determine whether',
        'verify', 'check if', 'confirm', 'validate'
    ])
    
    msgs = [
        {"role": "system", "content": "You are a semantic consistency checker. Your task is to determine if an answer correctly assesses and justifies the truth of a claim. Output only YES or NO followed by a brief explanation."},
        {"role": "user", "content": f"""QUESTION: {question}
ANSWER: {answer}

{'Does the ANSWER correctly assess the truth of the claim and justify it?' if truth_assessment else 'Does the ANSWER match what the question is asking for?'}
Answer YES or NO followed by a brief explanation."""}
    ]
    response = chat(msgs, temperature=0.1, max_tokens=64)
    is_consistent = response.strip().upper().startswith("YES")
    return is_consistent, response

def generate_candidates(question: str) -> List[Dict[str, str]]:
    """Generate candidates in parallel with rate limiting."""
    log.info(f"Generating {N_CANDIDATES} candidates")
    
    msgs = [
        SYSTEM_PROMPTS["gen"],
        {"role": "user", "content": GEN_USER.format(question=question)}
    ]
    
    def generate_one(i: int) -> Optional[Dict[str, str]]:
        """Generate a single candidate."""
        thread_id = f"gen-{i+1}"
        log.debug(f"[{thread_id}] Starting candidate generation")
        try:
            log.debug(f"[{thread_id}] Making chat request")
            cand_txt = chat(msgs, temperature=0.2, max_tokens=256)  # Lower temperature for better format adherence
            
            log.debug(f"[{thread_id}] Parsing response")
            cand = parse_candidate(cand_txt)
            if not cand:
                log.warning(f"[{thread_id}] Candidate failed to follow format. Raw response:\n{cand_txt.replace('\n', '\\n')}")
                return None
                
            # Check semantic consistency
            is_consistent, explanation = check_semantic_consistency(question, cand["answer"])
            if not is_consistent:
                log.warning(f"[{thread_id}] Candidate failed semantic consistency check: {explanation}")
                return None
                
            log.debug(f"[{thread_id}] Candidate parsed successfully")
            return cand
        except Exception as e:
            log.error(f"[{thread_id}] Failed to generate candidate: {str(e)}")
            return None
    
    # Generate candidates in parallel with a smaller number of workers
    log.info(f"Starting parallel generation with {min(3, MAX_WORKERS)} workers")
    with ThreadPoolExecutor(max_workers=min(3, MAX_WORKERS)) as executor:
        futures = [executor.submit(generate_one, i) for i in range(N_CANDIDATES)]
        candidates = [f.result() for f in futures if f.result() is not None]
    
    log.info(f"Generated {len(candidates)} valid candidates")
    return candidates

def critique(cand: Dict[str, str]) -> str:
    """Critique a candidate and return normalized response."""
    log.debug("Starting critique of candidate")
    msgs = [
        SYSTEM_PROMPTS["crit"],
        {"role": "user",
         "content": CRIT_USER.format(thought=cand['thought'], answer=cand['answer'])}
    ]
    return chat(msgs, temperature=0.2, max_tokens=256)  # Increased from 48

def process_candidates(cands: List[Dict[str, str]]) -> List[Dict[str, str]]:
    """Process candidates in parallel with rate limiting."""
    log.info(f"Processing {len(cands)} candidates")
    
    def process_one(cand: Dict[str, str], i: int) -> Dict[str, str]:
        """Process a single candidate."""
        thread_id = f"proc-{i+1}"
        log.debug(f"[{thread_id}] Starting candidate processing")
        try:
            log.debug(f"[{thread_id}] Starting critique")
            crit = critique(cand)
            
            if normalize_critic(crit) != "NOISSUE" and MAX_REFINE_ROUNDS:
                log.debug(f"[{thread_id}] Issues found, starting refinement")
                refined = refine(cand, crit)
                if refined:
                    log.debug(f"[{thread_id}] Refinement successful")
                else:
                    log.warning(f"[{thread_id}] Refinement failed, using original")
                return refined or cand
            else:
                log.debug(f"[{thread_id}] No issues found")
                return cand
        except Exception as e:
            log.error(f"[{thread_id}] Failed to process candidate: {str(e)}")
            return cand
    
    # Process candidates in parallel with a smaller number of workers
    log.info(f"Starting parallel processing with {min(3, MAX_WORKERS)} workers")
    with ThreadPoolExecutor(max_workers=min(3, MAX_WORKERS)) as executor:
        futures = [executor.submit(process_one, cand, i) for i, cand in enumerate(cands)]
        processed = [f.result() for f in futures]
    
    log.info(f"Processed {len(processed)} candidates")
    return processed

def refine(cand: Dict[str, str], critic: str) -> Dict[str, str]:
    """Refine a candidate based on critique."""
    log.debug("Starting refinement of candidate")
    
    # Dynamic template for refinement
    template = textwrap.dedent("""
        THOUGHT:
        <step-by-step mathematical reasoning>

        ANSWER:
        <one-line mathematical conclusion>
    """).strip()
    
    msgs = [
        SYSTEM_PROMPTS["ref"],
        {"role": "user",
         "content": f"YOUR PREVIOUS ANSWER:\n{cand['raw']}\n\nCRITIC SAYS:\n{critic}\n\nPlease rewrite using this template:\n\n{template}"}
    ]
    refined_txt = chat(msgs, temperature=0.4, max_tokens=256)  # Increased from 48
    new_cand = parse_candidate(refined_txt)
    if new_cand:
        log.debug("Refinement successful")
    else:
        log.warning("Refinement failed to parse, using original candidate")
    return new_cand or cand

def judge(cands: List[Dict[str, str]]) -> Tuple[int, List[str]]:
    """Judge candidates and extract rules."""
    log.debug(f"Starting judgment of {len(cands)} candidates")
    
    # Format candidates with letters
    cand_text = "\n\n".join(f"{chr(65+i)}. {c['raw']}" for i, c in enumerate(cands))
    
    msgs = [
        SYSTEM_PROMPTS["judge"],
        {"role": "user", "content": JUDGE_USER.format(candidates=cand_text)}
    ]
    verdict = chat(msgs, temperature=0.1, max_tokens=256)
    log.debug(f"Judge verdict: {verdict}")
    
    # Extract best candidate letter and rules
    best_match = re.search(r"Best:\s*([A-Z])", verdict)
    if not best_match:
        log.warning("No best candidate found in verdict, using first candidate")
        return 0, []
        
    best_letter = best_match.group(1)
    best_idx = ord(best_letter) - ord('A')
    
    # Validate index
    if best_idx < 0 or best_idx >= len(cands):
        log.warning(f"Invalid candidate index {best_idx}, using first candidate")
        return 0, []
    
    # Extract rules
    rules = []
    for line in verdict.split('\n'):
        if line.startswith('RULE:'):
            rules.append(line[5:].strip())
    
    return best_idx, rules

def thinking_loop(question: str) -> Tuple[Dict[str, str], List[str]]:
    """Main thinking loop with sequential processing."""
    log.debug("Starting thinking loop")
    cands = generate_candidates(question)
    
    # Handle case where no candidates were generated
    if not cands:
        log.warning("No valid candidates survived semantic check. Falling back to first parsed candidate.")
        # Try to get at least one candidate without semantic check
        msgs = [
            SYSTEM_PROMPTS["gen"],
            {"role": "user", "content": GEN_USER.format(question=question)}
        ]
        try:
            cand_txt = chat(msgs, temperature=0.2, max_tokens=256)
            fallback_cand = parse_candidate(cand_txt)
            if fallback_cand:
                return fallback_cand, []
        except Exception as e:
            log.error(f"Fallback generation failed: {str(e)}")
        
        # If even fallback fails, return error message
        return {"thought": "Error: No valid candidates generated.",
                "answer": "Please try again with a different prompt.",
                "raw": "THOUGHT:\nError: No valid candidates generated.\n\nANSWER:\nPlease try again with a different prompt."}, []
    
    # Process candidates sequentially
    cands = process_candidates(cands)
    
    best_idx, rules = judge(cands)
    log.debug(f"Selected candidate {best_idx+1} as best answer")
    return cands[best_idx], rules

# ──────────────────────────────────────────────────────────────────────
# 4  CLI
# ──────────────────────────────────────────────────────────────────────
if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.exit("Usage: python thinker_http_txt.py \"Your prompt here\"")
    prompt = sys.argv[1]

    t0 = time.time()
    best, rules = thinking_loop(prompt)
    dt = time.time() - t0

    print("\n=== FINAL ANSWER ===")
    print(best["raw"])  # Print the full THOUGHT + ANSWER format
    print("\n=== RULES EXTRACTED ===")
    print("\n".join(rules) or "(none)")
    print(f"\n(total time: {dt:.1f} s)")
