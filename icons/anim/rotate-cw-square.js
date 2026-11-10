import { animate, τmid } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="m9 8"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmid });
}
