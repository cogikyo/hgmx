import { animate, τmax } from "../motion.js";

export default function (el) {
  const core = el.querySelector('path[d^="M12 6.53"]');
  animate(core, { scaleX: [1, 0.85, 1.1, 1] }, { duration: τmax });
}
