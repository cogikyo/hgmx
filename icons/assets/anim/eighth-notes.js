import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const beam = el.querySelector("path");
  animate(beam, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const heads = el.querySelectorAll("circle");
  animate(heads, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
