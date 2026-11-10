import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bar = el.querySelector("line");
  animate(bar, { scaleX: [0, 1.25, 1] }, { duration: τmid });

  const dots = [el.querySelector('circle[cy="6"]'), el.querySelector('circle[cy="18"]')];
  animate(dots, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
