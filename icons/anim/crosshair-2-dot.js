import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const scope = el.querySelectorAll("circle, path");
  animate(scope, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
  const dot = el.querySelector('path[d="M12 12h.01"]');
  animate(dot, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
