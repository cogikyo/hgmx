import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arc = el.querySelector('path[d^="M3 12"]');
  animate(arc, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const head = el.querySelector('path[d^="M3 3"]');
  animate(head, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
  const key = el.querySelector("g");
  animate(key, { rotateY: [0, 360] }, { duration: τmid, delay: τmin });
}
