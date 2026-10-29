import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });
  const bars = el.querySelectorAll("line");
  animate(bars, { scaleY: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
