import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d="M3 9h18"]');
  animate(top, { translateY: [-6, 1, 0] }, { duration: τmid });

  const left = el.querySelector('path[d="M9 21V9"]');
  animate(left, { translateX: [-6, 1, 0] }, { duration: τmid, delay: τmin });
}
