import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const right = el.querySelector('path[d="M15 3v18"]');
  animate(right, { translateX: [6, -1, 0] }, { duration: τmid });

  const bottom = el.querySelector('path[d="M3 15h12"]');
  animate(bottom, { translateY: [6, -1, 0] }, { duration: τmid, delay: τmin });
}
