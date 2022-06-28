#include "wrapper.h"
#include <aubio/aubio.h>

void fvec_set_buffer(fvec_t *fvec, smpl_t *buf) {
    for (uint_t i = 0; i < sizeof(buf); i++) {
        fvec_set_sample(fvec, buf[i], i);
    }
}