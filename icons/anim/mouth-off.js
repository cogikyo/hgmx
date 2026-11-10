import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const mouth = el.querySelector("g");
  animate(mouth, { scaleY: [1, 1.15, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
