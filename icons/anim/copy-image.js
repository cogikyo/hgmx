import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { translateX: [-6, 1, 0], translateY: [-6, 1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M4 16"]');
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const hill = el.querySelector('path[d^="m13.4"]');
  animate(hill, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const sun = el.querySelector("circle");
  animate(sun, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
