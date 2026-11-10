import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector("g");
  animate(head, { scale: [1, 0.9, 1.05, 1] }, { duration: τmax });
  const eyes = el.querySelectorAll('path[d$="v-.5"]');
  animate(eyes, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
  const left = el.querySelector('path[d="M5 18H2"]');
  animate(left, { rotate: [0, 30, -10, 0] }, { duration: τmid, delay: τmin });
  const right = el.querySelector('path[d="M19 18h3"]');
  animate(right, { rotate: [0, -30, 10, 0] }, { duration: τmid, delay: τmin });
}
