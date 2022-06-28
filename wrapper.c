#include "wrapper.h"
#include <aubio/aubio.h>

// set all the values of fvec from an array of float32
void fvec_set_buffer(fvec_t *fvec, smpl_t *buf) {
    for (uint_t i = 0; i < fvec->length; i++) {
        fvec_set_sample(fvec, buf[i], i);
    }
}