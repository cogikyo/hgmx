import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M2 21"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const arrow = el.querySelector("g.arrow");
  animate(arrow, { translateX: [4, -2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
