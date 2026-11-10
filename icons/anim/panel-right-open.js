import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const divider = el.querySelector('path[d="M15 3v18"]');
  animate(divider, { translateX: [6, -1, 0] }, { duration: τmax });

  const chevron = el.querySelector('path[d^="m10 "]');
  animate(chevron, { translateX: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
