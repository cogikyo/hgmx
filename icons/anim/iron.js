import { animate, τmax } from "../motion.js";

export default function (el) {
  const iron = el.querySelectorAll('path:not([d="M3 22h10"])');
  animate(iron, { translateX: [0, -2, 1.5, 0] }, { duration: τmax });
}
