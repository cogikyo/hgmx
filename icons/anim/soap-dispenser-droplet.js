import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pump = el.querySelector("g");
  animate(pump, { translateY: [0, 1.5, 0] }, { duration: τmin });

  const drop = el.querySelector('path[d^="M19.29"]');
  animate(drop, { scaleY: [1, 0.85, 1.1, 1], scaleX: [1, 1.1, 0.95, 1] }, { duration: τmid, delay: τmin });
}
