import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d="M5 18H2"]');
  animate(left, { rotate: [0, 30, -10, 0] }, { duration: τmid, delay: τmin });
  const right = el.querySelector('path[d="M19 18h3"]');
  animate(right, { rotate: [0, -30, 10, 0] }, { duration: τmid, delay: τmin });
}
