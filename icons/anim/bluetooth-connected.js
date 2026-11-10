import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const rune = el.querySelector("path");
  animate(rune, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  animate(rune, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });

  const links = el.querySelectorAll("line");
  animate(links, { scaleX: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
