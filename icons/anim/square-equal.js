import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M7 10h10"]'), el.querySelector('path[d="M7 14h10"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
