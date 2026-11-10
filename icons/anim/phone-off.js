import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const handset = el.querySelector("g");
  animate(handset, { rotate: [0, -10, 6, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d^="M22 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
