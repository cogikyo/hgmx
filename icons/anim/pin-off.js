import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector("g");
  animate(head, { translateY: [-3, 0.5, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2 "]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
