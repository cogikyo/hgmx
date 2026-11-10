import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const kiwi = el.querySelector("circle");
  animate(kiwi, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const seeds = el.querySelectorAll("path");
  animate(seeds, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
