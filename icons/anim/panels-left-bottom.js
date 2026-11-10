import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d="M9 3v18"]');
  animate(left, { translateX: [-6, 1, 0] }, { duration: τmid });

  const bottom = el.querySelector('path[d="M9 15h12"]');
  animate(bottom, { translateY: [6, -1, 0] }, { duration: τmid, delay: τmin });
}
