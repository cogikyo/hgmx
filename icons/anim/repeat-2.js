import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const loop = [el.querySelector('path[d^="M13 18"]'), el.querySelector('path[d^="M11 6"]')];
  animate(loop, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin) });

  const up = el.querySelector('path[d^="m2 9"]');
  animate(up, { translateY: [4, -2, 0], opacity: [0, 1] }, { duration: τmid });

  const down = el.querySelector('path[d^="m22 15"]');
  animate(down, { translateY: [-4, 2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
