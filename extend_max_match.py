# run it as python3 extend_max_match.py in1.txt strings1.txt
# intervals are given as: [length of match, starting index in query, [starting indices in text]] 
#TODO: Assumption: Only 1 unique maximal match. But it could match in different locations in the text
#from numpy import array, sort

import sys

MAX_LENGTH = MAX_IDX = 0
actual_threshold = 32
THRESHOLD = actual_threshold >> 2
MATCH_SCORE = 5
INDEL_SCORE = -1
MISMATCH_SCORE = -4


def extend_match(AEM):
    # find the max match
    max_Q_strt_idx = AEM[MAX_IDX][1]
    max_Q_end_idx = max_Q_strt_idx + MAX_LENGTH - 1
    max_T_strt_idx = AEM[MAX_IDX][2]
    max_T_end_idx = max_T_strt_idx + MAX_LENGTH - 1

    last_accepted_index = MAX_IDX
    current_score = MAX_LENGTH

    # right extend till drop-off
    for i in range(MAX_IDX+1, len(AEM)):
    
        # if last accepted index is too high, break
        if i - last_accepted_index == THRESHOLD:
            break

        #find the closest exact match
        original_match_length = match_length = AEM[i][0]
        match_Q_idx =  AEM[i][1]
        T_indices = AEM[i][2:]
        add_Q_indels = add_T_indels = 0
        indels = mismatches = 0

        # distance between match segments in Q
        if match_Q_idx+match_length-1 <= max_Q_end_idx:
            # match lies with in already matched region
            continue
        elif match_Q_idx <= max_Q_end_idx:
            # overlapping matches, add indels in T match
            add_T_indels = max_Q_end_idx - match_Q_idx + 1
            
        # find the index closest to the right side of the match
        state = 0
        match_length -= add_T_indels
        for match_t in T_indices:
            if match_t + add_T_indels + match_length - 1 <= max_T_end_idx:
                # match lies with in already matched region
                continue
            elif match_t + add_T_indels <= max_T_end_idx:
                # overlapping matches, add indels in Q match
                add_Q_indels = max_T_end_idx - match_t - add_T_indels + 1
                match_length -= add_Q_indels
                state = 1
                break
            else:
                # not overlapping in text
                state = 2
                break

        if state == 0:
            # no new match info is found. Go to next
            continue
        if add_T_indels > 0:
            # indels in both T Q
            if add_Q_indels > 0:
                indels = add_Q_indels
                
            # indels in T only 
            else:
                indels = match_t + add_T_indels - max_T_end_idx

        else:
            # indels in Q only 
            if add_T_indels > 0:
                indels = match_Q_idx + add_Q_indels - max_Q_end_idx
            # no overlaps
            else:
                m1 = match_Q_idx - max_Q_end_idx - 1
                m2 = match_t - max_T_end_idx - 1
                mismatches = min(m1, m2)
                indels = max(m1, m2) - min(m1, m2)

        # if too many indels/mismatches, ignore
        score = mismatches*MISMATCH_SCORE + indels*INDEL_SCORE
        if score <= -THRESHOLD:
            continue
        else:
            # add the score to the existing match
            current_score += score + match_length*MATCH_SCORE
            max_Q_end_idx = match_Q_idx + original_match_length - 1
            max_T_end_idx = match_t + original_match_length - 1
            last_accepted_index = i

    

if __name__ == "__main__":
    num_lines = 0
    # create a static input array - avoid array.append()
    with open(sys.argv[1], 'r') as f:
        for i in f:
            num_lines += 1
        f.seek(0,0)
        inp_array = [[] for _ in range(num_lines)]
        
        # populate input array
        for idx,l in enumerate(f):
            temp = l.strip('\n').split()
            inp_array[idx] = [int(i) for i in temp]
            
            length_match = inp_array[idx][0]
            if length_match > MAX_LENGTH:
                MAX_LENGTH = length_match
                MAX_IDX = idx

        print(MAX_LENGTH, MAX_IDX)

        extend_match(inp_array)
