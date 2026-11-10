import { animate, τmax } from "../motion.js";

export default function (el) {
  const divider = el.querySelector('path[d="M9 3v18"]');
  animate(divider, { translateX: [-6, 1, 0] }, { duration: τmax });
}
