import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const divider = el.querySelector('path[d="M3 9h18"]');
  animate(divider, { translateY: [-6, 1, 0] }, { duration: τmax });

  const chevron = el.querySelector('path[d^="m15 "]');
  animate(chevron, { translateY: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
