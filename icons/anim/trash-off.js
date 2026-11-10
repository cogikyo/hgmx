import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lid = el.querySelector("g");
  animate(lid, { translateY: [0, -1.5, 0], rotate: [0, -8, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
