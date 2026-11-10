import { animate, τmax } from "../motion.js";

export default function (el) {
  const divider = el.querySelector('path[d="M3 9h18"]');
  animate(divider, { translateY: [-6, 1, 0] }, { duration: τmax });
}
