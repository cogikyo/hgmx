import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [6, -2, 0], opacity: [0, 1] }, { duration: τmax });

  const lights = el.querySelectorAll('path[d$="h.01"]');
  animate(lights, { scale: [0, 1.25, 1] }, { duration: τmin });
}
