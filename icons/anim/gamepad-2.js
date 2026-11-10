import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const pad = el.querySelector("g");
  animate(pad, { rotate: [0, -10, 6, 0] }, { duration: τmax });
  const buttons = el.querySelectorAll('line[x1="15"], line[x1="18"]');
  animate(buttons, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
