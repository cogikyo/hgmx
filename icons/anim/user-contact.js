import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("rect");
  animate(card, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const allPaths = el.querySelectorAll("path");
  const pins = [allPaths[0], allPaths[2]];
  animate(pins, { scale: [0, 1.75, 1] }, { duration: 0.3, delay: stagger(0.21) });

  const body = allPaths[1];
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0], opacity: [1] }, { duration: τmin, delay: τmid });
}
