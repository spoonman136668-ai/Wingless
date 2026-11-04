"""CNC-G9F: frozen Gaussian-Lloyd 2-bit expert representation, CPU only.

Four nonuniform fixed normalized centroids + per-row/block FP16 RMS. Not trained
on real Qwen weights. Decodes one block at a time, not the full matrix.
"""
import numpy as np

BLOCKS = (128,256)
LEVELS = np.array([-1.5104,-0.4528,0.4528,1.5104],dtype=np.float32)
CUTOFF = float((LEVELS[2]+LEVELS[3])/2)

def pack_companded_int2(weight, *, block=256):
    w=np.asarray(weight,dtype=np.float32)
    if (type(block) is not int or block not in BLOCKS or w.ndim!=2 or
        min(w.shape)<=0 or w.shape[1]%block or not np.isfinite(w).all()):
        raise ValueError("invalid companded INT2 expert")
    rows,cols=w.shape
    data=np.empty((rows,cols//4),dtype=np.uint8)
    scales=np.empty((rows,cols//block),dtype=np.float16)
    for k in range(cols//block):
        tile=w[:,k*block:(k+1)*block]
        rms=np.sqrt(np.mean(tile.astype(np.float64)**2,axis=1))
        s=rms.astype(np.float16)
        if not np.isfinite(s).all() or np.any((rms>0)&(s<=0)):
            raise ValueError("FP16 RMS scale underflow/overflow")
        scales[:,k]=s
        z=tile/np.where(s>0,s,np.float16(1)).astype(np.float32)[:,None]
        codes=np.where(z<-CUTOFF,0,np.where(z<0,1,
              np.where(z<CUTOFF,2,3))).astype(np.uint8)
        data[:,k*block//4:(k+1)*block//4]=(
            codes[:,0::4]|(codes[:,1::4]<<2)|
            (codes[:,2::4]<<4)|(codes[:,3::4]<<6))
    return data,scales

def _validate(data,scales,cols,block):
    if (not isinstance(data,np.ndarray) or data.dtype!=np.uint8 or data.ndim!=2 or
        not isinstance(scales,np.ndarray) or scales.dtype!=np.float16 or
        scales.ndim!=2 or type(cols) is not int or cols<=0 or
        type(block) is not int or block not in BLOCKS or cols%block or
        data.shape[0]<=0 or data.shape[1]*4!=cols or
        scales.shape!=(data.shape[0],cols//block) or
        not np.isfinite(scales).all() or np.any(scales<0)):
        raise ValueError("invalid companded INT2 representation")

def apply_companded_int2(inputs,data,scales,*,block=256):
    x=np.asarray(inputs,dtype=np.float32)
    if x.ndim!=2 or min(x.shape)<=0 or not np.isfinite(x).all():
        raise ValueError("invalid INT2 activations")
    _validate(data,scales,x.shape[1],block)
    out=np.zeros((x.shape[0],data.shape[0]),dtype=np.float32)
    for k in range(x.shape[1]//block):
        tile=data[:,k*block//4:(k+1)*block//4]
        codes=np.empty((tile.shape[0],block),dtype=np.uint8)
        for i in range(4): codes[:,i::4]=(tile>>(2*i))&3
        out+=(x[:,k*block:(k+1)*block]@LEVELS[codes].T)*scales[:,k][None,:].astype(np.float32)
    if not np.isfinite(out).all(): raise ValueError("nonfinite INT2 result")
    return out

def stored_bytes(data,scales,cols,*,block=256):
    _validate(data,scales,cols,block)
    return int(data.nbytes+scales.nbytes+24)
