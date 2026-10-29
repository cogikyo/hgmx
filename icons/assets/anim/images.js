import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const back = el.querySelector('path[d^="M4 8"]');
  animate(back, { translateX: [6, -1, 0], translateY: [-6, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const ridge = el.querySelector('path[d^="m22 11"]');
  animate(ridge, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const sun = el.querySelector("circle");
  animate(sun, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
