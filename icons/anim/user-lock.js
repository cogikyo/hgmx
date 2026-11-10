import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M9.5 15"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const shackle = el.querySelector('path[d^="M19 16"]');
  animate(shackle, { translateY: [0, 1.5, 0] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
