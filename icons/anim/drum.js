import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d="m2 2 8 8"]');
  animate(left, { rotate: [0, -30, 8, 0] }, { duration: τmid });
  const right = el.querySelector('path[d="m22 2-8 8"]');
  animate(right, { rotate: [0, 30, -8, 0] }, { duration: τmid, delay: τmin });
}
