import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const reels = el.querySelectorAll("circle");
  animate(reels, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid, delay: stagger(0.12) });

  const tape = el.querySelector("line");
  animate(tape, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
