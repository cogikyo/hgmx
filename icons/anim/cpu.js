import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const core = el.querySelector('rect[x="8"]');
  animate(core, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });

  const pins = el.querySelectorAll("g");
  animate(pins, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
