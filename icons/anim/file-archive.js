import { animate, stagger, τmax, τmid } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M13.659"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const teeth = el.querySelectorAll('path[d^="M8 7"], path[d^="M8 12"], path[d^="M8 18"]');
  animate([teeth[2], teeth[0], teeth[1]], { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const pull = el.querySelector("circle");
  animate(pull, { translateY: [-12, 1, 0] }, { duration: τmax });
}
