import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M4 12V4"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const digits = [el.querySelector("rect"), el.querySelector("g")];
  animate(digits, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
