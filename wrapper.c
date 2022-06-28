#include "wrapper.h"
#include <aubio/aubio.h>

// set all the values of fvec from an array of float32
void fvec_set_buffer(fvec_t *s, smpl_t* buf) {
  uint_t i;
  for ( i = 0; i < s->length; i++ )
  {
    s->data[i] = buf[i];
  }
}