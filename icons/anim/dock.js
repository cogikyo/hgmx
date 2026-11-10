import { animate, τmax } from "../motion.js";

export default function (el) {
  const dock = el.querySelector('path[d="M6 16h12"]');
  animate(dock, { translateY: [2, -1, 0], scaleX: [0.6, 1.1, 1] }, { duration: τmax });
}
