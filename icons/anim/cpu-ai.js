import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const pins = el.querySelectorAll("g:not(.ai)");
  animate(pins, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const ai = el.querySelector("g.ai");
  animate(ai, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
