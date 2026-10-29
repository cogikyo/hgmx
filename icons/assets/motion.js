export const τmax = 0.42;
export const τmin = 0.15;
export const τmid = τmax - τmin;

let started = null;

export function animate(...args) {
  const controls = Motion.animate(...args);
  started?.push(controls);
  return controls;
}

export const stagger = (...args) => Motion.stagger(...args);

export function record(play, el) {
  started = [];
  try {
    play(el);
    return started;
  } finally {
    started = null;
  }
}
