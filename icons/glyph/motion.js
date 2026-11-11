export const τmax = 0.42;
export const τmin = 0.15;
export const τmid = τmax - τmin;

let started = null;

export function animate(targets, keyframes, options) {
  // Motion's default gap repeats the dash at the path end, leaving a round-cap dot.
  if (keyframes.pathLength) keyframes = { pathSpacing: 1, ...keyframes };
  const controls = Motion.animate(targets, keyframes, options);
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
