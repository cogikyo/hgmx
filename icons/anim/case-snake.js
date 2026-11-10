import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const type = el.querySelector("g");
  animate(type, { scaleY: [1, 0.85, 1.15, 1] }, { duration: τmax });

  const separator = el.querySelector('path[d="M11 16h2"]');
  animate(separator, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
