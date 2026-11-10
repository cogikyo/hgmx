import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wrap = el.querySelector('path[d^="M3 12"]');
  animate(wrap, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const head = el.querySelector('path[d^="m16 16"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
