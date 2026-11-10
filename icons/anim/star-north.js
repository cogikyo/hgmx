import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const cross = el.querySelector("g");
  animate(cross, { scale: [0.5, 1.25, 1], rotate: [0, 90] }, { duration: τmax });

  const glints = el.querySelectorAll(":scope > path");
  animate(glints, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
