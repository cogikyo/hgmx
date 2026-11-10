import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bandage = el.querySelector("g");
  animate(bandage, { scaleX: [1, 1.1, 0.95, 1] }, { duration: τmax });

  const holes = el.querySelectorAll('path[d$="h.01"]');
  animate(holes, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
