import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });
  const dot = el.querySelector('path[d="M12 12h.01"]');
  animate(dot, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
