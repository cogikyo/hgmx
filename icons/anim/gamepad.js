import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const pad = el.querySelectorAll('line[x1="6"], line[x1="8"]');
  animate(pad, { scale: [1, 0.7, 1] }, { duration: τmid });
  const buttons = el.querySelectorAll('line[x1="15"], line[x1="18"]');
  animate(buttons, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
