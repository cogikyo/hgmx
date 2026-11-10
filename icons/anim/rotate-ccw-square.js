import { animate, τmid } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="m15 2-3"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmid });
}
